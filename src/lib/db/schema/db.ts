export interface MessDB {
	id: number;
	name: string;
}

export interface MessCapacityDB {
	id: number;
	mess_id: number;
	total_capacity: number;
	from_month: number;
	from_year: number;
}

export interface RegistrantDB {
	id: number;
	email: string;
	mess_id: number;
	month: number;
	year: number;
}
