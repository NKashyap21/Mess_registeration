import { error, json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { sql } from '../../../hooks.server';

export const GET: RequestHandler = async ({ url }) => {
	const email = url.searchParams.get('email') ?? '';
	if (email == '') {
		error(400, { message: 'email missing' });
	}
	const data = await sql`
		select
			mess_id
		from 
			registrants
		where 
			email = ${email}
	`;
	if (data.at(0) == undefined) {
		return json({ mess_id: -1 });
	}
	return json({ mess_id: data.at(0)!['mess_id'] });
};
