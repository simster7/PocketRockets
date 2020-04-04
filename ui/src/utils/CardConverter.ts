let SUIT_MAP: { [id: number]: string } = {
	0: 'S',
	1: 'H',
	2: 'C',
	3: 'D'
};

let RANK_MAP: { [id: number]: string } = {
	0: '2',
	1: '3',
	2: '4',
	3: '5',
	4: '6',
	5: '7',
	6: '8',
	7: '9',
	8: 'T',
	9: 'J',
	10: 'Q',
	11: 'K',
	12: 'A'
};

export function cardIdToString(cardId: number): string {
	return RANK_MAP[cardId % 13] + SUIT_MAP[Math.floor(cardId / 13)];
}
