import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ parent }) => {
	const par = await parent();
	if (par.user['user_type'] != 2) {
		if (par.user['user_type'] == 1) {
			throw redirect(307, '/mess');
		} else {
			throw redirect(307, '/student');
		}
	}
};
