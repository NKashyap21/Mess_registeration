import { PRIVATE_API_URL } from '$env/static/private';
import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
import { resolve } from '$app/paths';

export const load: LayoutServerLoad = async ({ fetch, url }) => {
	if (url.pathname != resolve('/login')) {
		try {
			const res = await fetch(PRIVATE_API_URL + '/getUser', {
				method: 'GET',
				credentials: 'include'
			});
			if (res.status != 200) {
				await fetch(PRIVATE_API_URL + '/logout', { method: 'POST', credentials: 'include' });
				throw redirect(307, resolve('/login'));
			}
			let userData = await res.json();
			userData = userData['data'];
			return {
				user: userData
			};
		} catch (e) {
			console.error('Failed layout server load at /');
			console.error(e);
			throw redirect(307, resolve('/login'));
		}
	}
};
