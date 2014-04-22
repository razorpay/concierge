<?php

/*
|--------------------------------------------------------------------------
| Application Routes
|--------------------------------------------------------------------------
|
| Here is where you can register all of the routes for an application.
| It's a breeze. Simply tell Laravel the URIs it should respond to
| and give it the Closure to execute when that URI is requested.
|
*/

Route::group(array('before' => 'auth'), function()
{
Route::get('/groups', 'HomeController@getGroups');

Route::get('/manage/{group_id}', 'HomeController@getManage');

Route::post('/manage/{group_id}', array('before'=>'csrf', 'uses'=>'HomeController@postManage'));

Route::post('/manage/{group_id}/terminate', array('before'=>'csrf', 'uses'=>'HomeController@postTerminate'));

Route::get('/logout', 'HomeController@getLogout');

Route::get('/password', 'HomeController@getPassword');

Route::post('/password', array('before'=>'csrf', 'uses'=>'HomeController@postPassword'));

Route::get('/users', array('before'=>'admin', 'uses'=>'HomeController@getUsers'));

Route::post('/users', array('before'=>'admin|csrf', 'uses'=>'HomeController@postUsers'));
});

Route::group(array('before' => 'guest'), function()
{
Route::get('/', 'HomeController@getIndex');

Route::post('/signin', array('before'=>'csrf', 'uses'=>'HomeController@postSignin'));

Route::post('/duologin', 'HomeController@postDuologin');

});

Route::get('/invite/{token}', 'HomeController@getInvite');