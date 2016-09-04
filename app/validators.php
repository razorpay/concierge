<?php
/*
* app/validators.php
* Custom Validation Rules
*/

Validator::extend('alpha_spaces', function($attribute, $value)
{
    return preg_match('/^[\pL\s]+$/u', $value);
});

Validator::extend('razorpay_email', function ($attribute, $value)
{
    return checkEmailDomain($value);
});
