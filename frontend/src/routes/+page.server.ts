import { PUBLIC_API_URL } from '$env/static/public';
import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	try {
		const res = await fetch(PUBLIC_API_URL + '/getUser', { method: 'GET', credentials: 'include' });
		if (res.status != 200) {
			fetch(PUBLIC_API_URL + '/logout', { method: 'POST', credentials: 'include' });
			redirect(301, '/login');
		}
		let userData = await res.json();
		userData = userData['data'];
		if (userData['user_type'] == 0) {
			throw redirect(307, '/student');
		} else if (userData['user_type'] == 1) {
			throw redirect(307, '/mess');
		} else if (userData['user_type'] == 2) {
			throw redirect(307, '/admin');
		} else {
			fetch(PUBLIC_API_URL + '/logout', { method: 'POST', credentials: 'include' });
			redirect(301, '/login');
		}
	} catch (e) {
		console.error(e);
		redirect(301, '/login');
	}
};
