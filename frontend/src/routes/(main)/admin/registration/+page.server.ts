import { PRIVATE_API_URL } from '$env/static/private';
import type { PageServerLoad } from './$types';

export const ssr = false;

export const load: PageServerLoad = async ({ fetch }) => {
	return {
		registrationState: await (
			await fetch(PRIVATE_API_URL + '/students/isRegistrationOpen', {
				method: 'get',
				credentials: 'include'
			})
		).json(),
		messStats: await (
			await fetch(PRIVATE_API_URL + '/office/messStatsGrouped', {
				method: 'GET',
				credentials: 'include'
			})
		).json()
	};
};
