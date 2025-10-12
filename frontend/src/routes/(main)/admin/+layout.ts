import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = ({ url }) => {
	console.log(url.pathname);
	if (url.pathname == '/admin') {
		redirect(303, '/admin/registration');
	}
};
