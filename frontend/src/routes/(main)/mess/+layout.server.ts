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
	if (userData['user_type'] != 1) {
		if (userData['user_type'] == 2) {
			throw redirect(307, '/admin');
		} else {
			throw redirect(307, '/student');
		}
	}
};
