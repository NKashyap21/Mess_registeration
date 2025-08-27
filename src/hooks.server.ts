export { handle } from '$lib/auth';

import { DB_URL, REDIS_URL } from '$env/static/private';
import type { RegistrationData } from '$lib/types/interface';
import postgres from 'postgres';
import { createClient } from 'redis';

export const sql = postgres(DB_URL);

export const redisClient = createClient({ url: REDIS_URL });
redisClient.connect();

async function redisQueueWorker() {
	const blockClient = createClient({ url: REDIS_URL });
	blockClient.connect();

	console.log('Redis Queue Worker Started to insert registrants.');

	while (true) {
		const data = JSON.parse((await blockClient.blPop('registration:registrants', 0)).element);
		console.log(data);
		try {
			console.log(`Registrant Email: ${data.email}  Mess: ${data.messId}`);
			await sql`
			INSERT INTO 
				registrants
				(email, mess_id, month, year)
			VALUES
				(${data.email}, ${data.messId}, ${new Date().getMonth() + 1}, ${new Date().getFullYear()})
			`;
			console.log('Registrant inserted');
		} catch (err) {
			console.error(
				`Failed to add registrant ${data.email} to mess ${data.messId} due to error ${err}`
			);
			await new Promise((r) => {
				setTimeout(r, 500);
			});
		}
	}
}

redisQueueWorker();
