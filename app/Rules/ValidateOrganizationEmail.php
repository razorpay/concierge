<?php

namespace App\Rules;

use Illuminate\Contracts\Validation\Rule;

class ValidateOrganizationEmail implements Rule
{
    /**
     * Determine if the validation rule passes.
     *
     * @param  string  $attribute
     * @param  mixed  $value
     * @return bool
     */
    public function passes($attribute, $value)
    {
        $domain = config('concierge.google_domain');
        $email_parts = explode('@', $value);
        $email_domain = end($email_parts);
        return ($email_domain === $domain);
    }

    /**
     * Get the validation error message.
     *
     * @return string
     */
    public function message()
    {
        return ':attribute should be an organization email address';
    }
}
