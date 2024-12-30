import { authenticateUser } from '$lib/server/auth';
import { redirect, type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	//log.bold(`📣 NEW REQUEST IS BEING MADE FROM ${event.url.pathname}`);
	//log.hooks('hooks.server.ts');
	console.log('Hook', event.url.pathname);
	event.locals.user = await authenticateUser(event);

	console.log(event.url.pathname, event.locals.user);

	if (event.url.pathname == '/')
		if (event.locals.user == null) return await resolve(event);
		else return redirect(301, '/archive');

	if (!event.locals.user) {
		redirect(301, '/');
	}

	const response = await resolve(event);

	return response;
};
