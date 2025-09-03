export interface MessData {
	id: number;
	name: string;
	total_registrants: number;
	total_capacity: number;
}

export interface RegistrationData {
	email: string;
	messId: number;
}

export interface SearchEmailRequest {
	email: string;
}

export interface SearchEmailResponse {
	mess_id: number;
}
