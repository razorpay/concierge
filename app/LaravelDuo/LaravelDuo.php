<?php 
namespace LaravelDuo;
class LaravelDuo extends Duo
{

    private $_AKEY;
    private $_IKEY;
    private $_SKEY;
    private $_HOST;

    public function __construct() { 
        $_AKEY = \Config::get('custom_config.duo_akey');
        $_IKEY = \Config::get('custom_config.duo_ikey');
        $_SKEY = \Config::get('custom_config.duo_skey');
        $_HOST = \Config::get('custom_config.duo_host');
    }

    public function get_akey()
    {
        return $this->_AKEY;
    }

    public function get_ikey()
    {
        return $this->_IKEY;
    }

    public function get_skey()
    {
        return $this->_SKEY;
    }

    public function get_host()
    {
        return $this->_HOST;
    }

} 