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

        $users = array(
            array(
                'name' => 'Shashank Kumar',
                'username' => 'shk',
                'password' => Hash::make('password'),
                'created_at' => date('Y-m-d H:i:s'),
                'updated_at' => date('Y-m-d H:i:s'),
                'admin' => TRUE
            ),
            array(
                'name' => 'Harshil Mathur',
                'username' => 'harshil',
                'password' => Hash::make('password'),
                'created_at' => date('Y-m-d H:i:s'),
                'updated_at' => date('Y-m-d H:i:s'),
                'admin' => TRUE
            ),
            array(
                'name' => 'Abhay Bir Singh Rana',
                'username' => 'nemo',
                'password' => Hash::make('password'),
                'created_at' => date('Y-m-d H:i:s'),
                'updated_at' => date('Y-m-d H:i:s'),
                'admin' => FALSE
            )
        );
        DB::table('users')->insert($users);

	}

}