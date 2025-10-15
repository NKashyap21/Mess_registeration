import { PUBLIC_API_URL } from '$env/static/public';
import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch }) => {
	const res = await fetch(PUBLIC_API_URL + '/getUser', { method: 'GET', credentials: 'include' });
	if (res.status != 200) {
		fetch(PUBLIC_API_URL + '/logout', { method: 'POST', credentials: 'include' });
		redirect(301, '/login');
	}
	let userData = await res.json();
	userData = userData['data'];
	if (userData['user_type'] != 0) {
		if (userData['user_type'] == 1) {
			throw redirect(307, '/mess');
		} else {
			throw redirect(307, '/admin');
		}
	}
	return {
		user: userData,
		regData: await (
			await fetch(PUBLIC_API_URL + '/students/isRegistrationOpen', {
				method: 'GET',
				credentials: 'include'
			})
		).json(),
		userMessData: await (
			await fetch(PUBLIC_API_URL + '/students/getMess', { method: 'GET', credentials: 'include' })
		).json()
	};
};
