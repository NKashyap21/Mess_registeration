import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = ({ url }) => {
	if (url.pathname == '/admin') {
		redirect(303, '/admin/registration');
	}
};
