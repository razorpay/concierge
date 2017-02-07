<?php

function checkEmailDomain($email, $domain = null)
{
    $domain = $domain ?: env('COMPANY_DOMAIN');

    $email_parts = explode('@', $email);

    $email_domain = end($email_parts);

    return $email_domain === $domain;
}
