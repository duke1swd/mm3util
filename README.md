# mm3util
A utility to manipulate mailman3 mailing lists

Can add or delete users
Can get the list of users
Can subscribe and unsubscribe users to lists
Can get the list of users who are subscribed to a list.`


Examples
mm3util domains
	returns the list of all domains
mm3util lists
	returns the name of all lists in the default domain.

mm3util user add <preferred email address> 'Display Name' <option1> <value1> ...  
mm3util user delete <perferred email address>
mm3util user show
	returns the entire user record

mm3util subscribe <list> <preferred email address> <option1> <value1> ...  
mm3util unsubscribe <list> <perferred email address>

mm3util list <list>
	returns all emails on the list


Configuration is managed by a json file.
Default configuration is /opt/mailman/mm/mm3util.cfg
File looks like this:
	{
		url: 'http://localhost:8001/3.1/',
		username: 'restuser',
		password: 'restpass',
	}

Non standard config file specified with '-c' flag
Select domain for all commands but "domains" with '-d' flag


How this is used by the web site:
When a user subscribes to a list, we first attempt to create a user with no display name and that email address.
THen we subscribe the email address to the list.

Unsubscribing an address works as expected.

mm3util clean is a utility that finds all addresses that are not subscribed to a list deletes them, then finds all users with no addresses and deletes those.

reconcilliation is performed by using mm3util list show <listname>, which shows all email addresses in that list.
