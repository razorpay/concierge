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

Route::post('/manage/{group_id}', 'HomeController@postManage');

Route::get('/logout', 'HomeController@getLogout');
});

Route::group(array('before' => 'guest'), function()
{
Route::get('/', 'HomeController@getIndex');

Route::post('/signin', 'HomeController@postSignin');

Route::post('/duologin', 'HomeController@postDuologin');
});