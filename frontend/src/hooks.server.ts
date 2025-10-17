import { PRIVATE_API_URL } from '$env/static/private';
import { redirect, type Handle, type HandleFetch } from '@sveltejs/kit';

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

export const handleFetch: HandleFetch = async ({ request, fetch, event }) => {
	if (request.url.startsWith(PRIVATE_API_URL)) {
		request.headers.set('cookie', event.request.headers.get('cookie') ?? '');
	}
	return fetch(request);
};
