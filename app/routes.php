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

//Routes for logged in user
Route::group(array('before' => 'auth'), function()
{
	Route::get('/groups', 'HomeController@getGroups');

	Route::get('/manage/{group_id}', 'HomeController@getManage');

	Route::post('/manage/{group_id}', array('before'=>'csrf', 'uses'=>'HomeController@postManage'));

	Route::post('/manage/{group_id}/terminate', array('before'=>'csrf', 'uses'=>'HomeController@postTerminate'));

	Route::get('/logout', 'HomeController@getLogout');

	Route::get('/password', 'HomeController@getPassword');

	Route::post('/password', array('before'=>'csrf', 'uses'=>'HomeController@postPassword'));
});

//Routes for site admin
Route::group(array('before' => 'auth|admin'), function()
{
	Route::get('/users', 'HomeController@getUsers');

	Route::post('/users', array('before'=>'csrf', 'uses'=>'HomeController@postUsers'));

	Route::get('/users/add', 'HomeController@getAddUser');

	Route::post('/users/add', array('before'=>'csrf', 'uses'=> 'HomeController@postAddUser'));

});

//Routes for non-logged in user
Route::group(array('before' => 'guest'), function()
{
	Route::get('/', 'HomeController@getIndex');

	Route::post('/signin', array('before'=>'csrf', 'uses'=>'HomeController@postSignin'));

	Route::post('/duologin', 'HomeController@postDuologin');

});

Route::get('/invite/{token}', 'HomeController@getInvite');