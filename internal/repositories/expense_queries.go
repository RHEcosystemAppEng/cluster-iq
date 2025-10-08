package repositories

const (
	ExpensesTable = "expenses"

	// SelectExpensesQuery returns every expense in the inventory ordered by instanceID
	SelectExpensesQuery = `SELECT * FROM expenses`
	// InsertExpensesQuery inserts into a new expense for an instance
	InsertExpensesQuery = `
		INSERT INTO expenses (
			instance_id,
			date,
			amount
		) VALUES (
			(SELECT id FROM instances WHERE instance_id=:instance_id),
			:date,
			:amount
		) ON CONFLICT (instance_id, date) DO UPDATE SET
			amount = EXCLUDED.amount
	`
	// SelectExpensesByInstanceQuery returns expense in the inventory for a specific InstanceID
	SelectExpensesByInstanceQuery = `
		SELECT * FROM expenses
		WHERE instance_id = $1
		ORDER BY date
	`
	SelectLastExpensesQueryVOPT = `
		WITH ranked_expenses AS (
			SELECT
				instance_id,
				date,
				amount,
				ROW_NUMBER() OVER (PARTITION BY instance_id ORDER BY date DESC) AS rn
			FROM expenses
		)
		SELECT
			re.instance_id,
			re.date,
			re.amount
		FROM ranked_expenses re
		JOIN instances i ON re.instance_id = i.id
		WHERE re.rn = 1
		AND i.status != 'Terminated'
		AND re.date < '$1';
	`
	// SelectLastExpensesQuery returns the last expense for every instance older
	// than 1 day. This is used for obtaining the list of instances that need
	// Billing information update because all the instances returned by this
	// query doesn't have expenses for the current day
	SelectLastExpensesQuery = `
		SELECT
				instances.id
		FROM
				instances
		LEFT JOIN (
				SELECT
						instance_id,
						MAX(date) AS last_expense_date
				FROM
						expenses
				GROUP BY
						instance_id
		) AS last_expenses
		ON
				instances.id = last_expenses.instance_id
		WHERE
				last_expenses.last_expense_date IS NULL
				OR last_expenses.last_expense_date < CURRENT_DATE - INTERVAL '1 day';
	`
)
