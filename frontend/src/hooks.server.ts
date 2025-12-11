import { PRIVATE_API_URL } from '$env/static/private';
import { redirect, type Handle, type HandleFetch } from '@sveltejs/kit';

import { resolve as resolveUrl } from '$app/paths';

export const handle: Handle = async ({ event, resolve }) => {
	if (event.url.pathname == resolveUrl('/login')) {
		console.log('IN LOGIN HANDLER');
		return resolve(event);
	}
	console.log('HANDLER');
	const jwt = event.cookies.get('mess_jwt');
	console.log('JWT: ' + jwt);
	if (!jwt) {
		console.log('THIS');
		throw redirect(307, resolveUrl('/login'));
	}
	return resolve(event);
};

export const handleFetch: HandleFetch = async ({ request, fetch, event }) => {
	if (request.url.startsWith(PRIVATE_API_URL)) {
		request.headers.set('cookie', event.request.headers.get('cookie') ?? '');
	}
	return fetch(request);
};
