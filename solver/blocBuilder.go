package solver

import(
	"sync"

	TJ "vintz.fr/nonogram/tabjeu"
)

// buildAllSequences construit la liste des ensembles de lignes (ou colonnes)
// à partir d'une liste de listes de longueurs de blocs
func buildAllSequences(taille int, seqs []TJ.BlocCount) lineListSet {
	result := make(lineListSet, taille)

	ch := make(chan indexedLineSet)
	wg := new(sync.WaitGroup)

	// construit les ensembles de lignes possibles en parallèle
	for i := range seqs {
		wg.Add(1)
		go func(taille int, sc TJ.BlocCount, idx int) {
			lines := buildSequences(taille, sc)
			ch <- indexedLineSet{num: idx, lines: lines}
		}(taille, seqs[i], i)
	}

	// reçoit les resultats pour chaque ligne
	go func() {
		for lines := range ch {
			result[lines.num] = lines.lines
			wg.Done()
		}
	}()

	wg.Wait()
	return result
}

// buildSequences construit l'ensemble des lignes (ou colonnes) possibles
// à partir d'une liste de longueurs de blocs et de la taille de la ligne
func buildSequences(taille int, blocs TJ.BlocCount) lineList {
	// cas particulier des lignes complètement vides
	if len(blocs) == 0 {
		ligneVide := make(TJ.LigneJeu, taille)
		for i := range ligneVide {
			ligneVide[i] = cellVide
		}
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

	// essaye succesivement le bloc sur toutes les positions possibles
	for startPos := 0; startPos <= (taille - longMini); startPos++ {

		tailleSeqCourante := taille // dernier blocs
		if !lastBloc {
			//  non-derniers blocs
			tailleSeqCourante = startPos + blocs[0] + 1
		}
		seqCourante := make(TJ.LigneJeu, tailleSeqCourante)
		for i := 0; i < tailleSeqCourante; i++ {
			// place des cellules pleines entre startPos et la fin du bloc
			if (startPos) <= i && (i < startPos+blocs[0]) {
				seqCourante[i] = cellPlein
				continue
			}
			// place des cellules vides ailleurs (avant et/ou après le bloc)
			seqCourante[i] = cellVide
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
