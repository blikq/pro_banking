# pro_banking

A production ready banking solution

## BUGS:
	. Roles should be casted to lowercase
	. endpoint for normal customer auth
	. auth_service not checking for expiration (for developmet) (feature)

## RUNNING
Make a .env file in each service referencing from the .env.example files in each service.
Install dependencies

### make run :

	IN order:
		auth_service
		tx_service


## ENDPOINTS
### AUTH_SERVICE
		. /api/register (service)
		. /api/login (service)
		. /api/auth-user (service)

### TX_SERVICE
		. /api/create-tx (customer & service)
		. /api/get-tx (customer)