<?php

class DatabaseSeeder extends Seeder {

	/**
	 * Run the database seeds for setting up inital set of users
     * 
	 *
	 * @return void
	 */
	public function run()
	{
		Eloquent::unguard();
        $users=Config::get('custom_config.users');
        DB::table('users')->insert($users);

	}

}