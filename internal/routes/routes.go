package routes

// Routes-- #######################################################

// AUTHENTICATION Routes

const REGISTER = "/register-user"
const LOGIN = "/login-user"
const DRIVERREGISTER = "/register-driver"
const LOGINDRIVER = "/login-driver"
const LOGOUT = "/logout"

// REQUESTER Routes
const USER = "/user"
const CREATE_REQUEST_INTENT = "/user/request-intent" // When triggered notifies the most nearby Driver //

// Driver Routes
const DRIVER = "/driver"
const ACCPET_REQUEST_INTENT = "/driver/request-intent-accept"
