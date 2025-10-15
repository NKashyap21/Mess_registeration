import { redirect } from '@sveltejs/kit';
import { OAuth2Client } from 'google-auth-library';
import { GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET } from '$env/static/private';
import { goto } from '$app/navigation';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ url }) => {
	// const redirectURL = 'http://localhost:5173/oauth';
	const redirectURL = 'http://localhost:3000/oauth';
	const code = await url.searchParams.get('code');

	console.log('code', code);

	try {
		const oAuth2Client = new OAuth2Client(GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, redirectURL);

		const r = await oAuth2Client.getToken(code);
		oAuth2Client.setCredentials(r.tokens);
		console.log('Auth Tokens recieved');
		const user = oAuth2Client.credentials;
		console.log('Credentials:', user);

		// Send id_token to your backend
		const loginToBackend = async (idToken: string) => {
			const res = await fetch('http://localhost:8080/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ token: idToken })
			});

			if (!res.ok) {
				console.error('Backend login failed:', await res.text());
				return;
			}

			const data = await res.json();
			console.log('Backend login success:', data);

			// Store your backend JWT
			localStorage.setItem('jwt', data.data.token);

			// Redirect to student page
			// goto('/student');
		};

		await loginToBackend(user.id_token); // Call backend login
	} catch (err) {
		console.log('Error:', err);
	}
	throw redirect(303, '/student');
};
