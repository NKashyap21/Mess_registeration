import type { Actions } from './$types';

export const actions = {
	default: ({ request }) => {
		const data = request.formData();
	}
} satisfies Actions;
