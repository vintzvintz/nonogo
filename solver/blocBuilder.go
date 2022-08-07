package solver

import(
	"sync"

	TJ "vintz.fr/nonogram/tabjeu"
)

// buildAllSequences construit la liste des ensembles de lignes (ou colonnes)
// à partir d'une liste de listes de longueurs de blocs
func buildAllSequences(taille int, blocs []TJ.BlocCount) lineListSet {
	result := make(lineListSet, taille)

	ch := make(chan indexedLineSet)
	wg := new(sync.WaitGroup)

	// construit les ensembles de lignes possibles en parallèle
	for i := range blocs {
		wg.Add(1)
		go func(taille int, sc TJ.BlocCount, idx int) {
			lines := buildSequences(taille, sc)
			ch <- indexedLineSet{num: idx, lines: lines}
		}(taille, blocs[i], i)
	}

	// reçoit les resultats pour chaque ligne
	go func() {
		for lines := range ch {
			result[lines.num] = lines.lines
			wg.Done()
		}
	}()

	wg.Wait()
	close(ch)
	return result
}

// buildSequences construit recursivement toutes les combinaisons possibles
// correspondant à une liste de longueurs de blocs pour une ligne ou une colonne
func buildSequences(taille int, blocs TJ.BlocCount) lineList {
	// cas particulier des lignes complètement vides
	if len(blocs) == 0 {
		ligneVide := make(TJ.LigneJeu, taille) // zero-value = Vide 
//		for i := range ligneVide {
//			ligneVide[i] = cellVide
//		}
		return lineList{ligneVide}
	}

	// cas particulier pour le dernier bloc (pas de séparateur à la fin)
	lastBloc := len(blocs) == 1

	// calcule le nb de cellules mini pour placer tous les blocs avec un espacement de 1
	var longMini int = len(blocs) - 1 // nb de cellules vides intercalaires
	for _, s := range blocs {
		longMini += s
	}

	result := make(lineList, 0)

	// place succesivement le bloc sur toutes les positions possibles
	for startPos := 0; startPos <= (taille - longMini); startPos++ {

		tailleSeqCourante := taille // dernier blocs
		if !lastBloc {
			//  non-derniers blocs
			tailleSeqCourante = startPos + blocs[0] + 1
		}
		seqCourante := make(TJ.LigneJeu, tailleSeqCourante)
		for i := 0; i < tailleSeqCourante; i++ {
			// par défaut, toutes les cellules sont initialisées vides, inutile de les remplir explicitement
			//seqCourante[i] = cellVide

			// remplit la ligne à partir de startPos jusqu'à la fin du bloc
			if (startPos) <= i && (i < startPos+blocs[0]) {
				seqCourante[i] = TJ.Plein
				continue
			}
		}

		// Si c'est le dernier bloc on renvoie juste les séquences courantes
		if lastBloc {
			result = append(result, seqCourante)
		}

		// si ce n'est pas le dernier bloc, appel récursif pour les séquences restantes
		if !lastBloc {
			seqsSuivantes := buildSequences(taille-tailleSeqCourante, blocs[1:])

			//  concatène les séquences suivantes avec la séquence courante
			for i := range seqsSuivantes {
				seq := append(seqCourante, seqsSuivantes[i]...)
				result = append(result, seq)
			}
		}
	}
	return result
}
