package queueManager

const(
	SQL_GET_BALANCE_BY_USER_ID = `
		SELECT ac.balance FROM public.users as u
		JOIN public.accounts as ac
			ON ac.owner_id=u.user_id
		WHERE u.user_id=$1;
	`
	SQL_UPDATE_BALANCE_BY_USER_ID = `
		UPDATE public.accounts
		SET balance=$1
		WHERE owner_id=$2;
	`

	SQL_INSERT_TRANSACTION_HISTORY_ENTRY = `
		INSERT INTO public.transaction_history (transaction_id, account_id, operation, value, date)
		VALUES ($1, $2, $3, $4, $5);
	`

	SQL_SELECT_TRANSACTION_HISTORY_ENTRIES_BY_USER_ID = `
		SELECT th.account_id, th.operation, th.value, th.date
		FROM public.transaction_history as th
		JOIN public.accounts as ac
			ON ac.account_id=th.account_id
		WHERE ac.owner_id=$1
		ORDER BY th.date ASC;
	`

	SQL_SELECT_ACCOUNT_ID_BY_USER_ID = `
		SELECT account_id
		FROM public.accounts
		WHERE owner_id=$1;
	`
)