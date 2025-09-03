import type { MessData } from '$lib/types/interface';
import { redisClient, sql } from '../hooks.server';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async (event) => {
	const session = await event.locals.auth();
	return {
		session: session,
		messData: (await sql`
		SELECT DISTINCT ON (m.name) 
			m.id,
			m.name, 
			COALESCE(mc.total_capacity, 0)::INT as total_capacity,
			COALESCE(COUNT(r.id), 0)::INT as total_registrants
		FROM
			mess m
		LEFT JOIN
			mess_capacities mc
			ON mc.mess_id = m.id AND (mc.from_year, mc.from_month) <= (${new Date().getFullYear()}, ${new Date().getMonth() + 1} )
		LEFT JOIN
			registrants r		
			ON r.mess_id = m.id AND r.month = ${new Date().getMonth() + 1} AND r.year = ${new Date().getFullYear()}
		GROUP BY
			m.name, mc.total_capacity, mc.from_year, mc.from_month, m.id
		ORDER BY	
			m.name ASC, mc.from_year DESC, mc.from_month DESC NULLS LAST
		`) as MessData[],
		registrationLive: (await redisClient.get('mess:live')) == 'true',
		registeredMess: (await sql` 
		SELECT 
			m.name
		FROM 
			registrants r
		JOIN
			mess m
			ON r.mess_id = m.id
		WHERE 
			r.email = ${session?.user?.email ?? ''} AND r.month = ${new Date().getMonth() + 1} AND r.year = ${new Date().getFullYear()}
			`) as any[]
	};
};
