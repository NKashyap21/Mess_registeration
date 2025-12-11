import { redirect } from '@sveltejs/kit';
import type { LayoutLoad } from './$types';
import { resolve } from '$app/paths';

export const load: LayoutLoad = ({ url }) => {
	if (url.pathname == resolve('/admin')) {
		redirect(307, resolve('/admin/registration'));
	}
};
