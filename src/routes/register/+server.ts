import type { RegistrationData } from '$lib/types/interface';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { redisClient } from '../../hooks.server';

export const POST: RequestHandler = async ({ request }) => {
	try {
		const data = (await request.json()) as RegistrationData;
		await redisClient.rPush('mess:registrants', JSON.stringify(data));
		return json({ success: true, error: null });
	} catch (error) {
		return json({ success: false, error: error });
	}
};
