import { redirect, type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	if (event.url.pathname == '/login') {
		return resolve(event);
	}
	const jwt = event.cookies.get('jwt');
	if (!jwt) {
		throw redirect(307, '/login');
	}
	return resolve(event);
};
