<?php

function checkEmailDomain($email, $domain = null)
{
    $domain = $domain ?: $_ENV['company_domain'];

    $email_parts = explode('@', $email);

    $email_domain = end($email_parts);
    
    return $email_domain === $domain;
}
