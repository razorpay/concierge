<?php

use Illuminate\Console\Command;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Input\InputArgument;

class croncreator extends Command {

	/**
	 * The console command name.
	 *
	 * @var string
	 */
	protected $name = 'custom:croncreator';

	/**
	 * The console command description.
	 *
	 * @var string
	 */
	protected $description = 'Command description.';

	/**
	 * Create a new command instance.
	 *
	 * @return void
	 */
	public function __construct()
	{
		parent::__construct();
	}

	/**
	 * Execute the console command.
	 *
	 * @return mixed
	 */
	public function fire()
	{
		$output = shell_exec('crontab -l');
		$cron="* * * * *  ".$GLOBALS['php_path']." ".$GLOBALS['artisan_path']." custom:leasemanager";
		if(strpos($output,$cron)===FALSE)
		{ 
			file_put_contents('/tmp/crontab.txt', $output.$cron.PHP_EOL);
			exec('crontab /tmp/crontab.txt');
			$this->info("Cron Added for running every minute");
		}
		else
		{
			$this->info("Cron Already exists");
		}
		
	}

	/**
	 * Get the console command arguments.
	 *
	 * @return array
	 */
	/*protected function getArguments()
	{
		return array(
			array('example', InputArgument::REQUIRED, 'An example argument.'),
		);
	}*/

	/**
	 * Get the console command options.
	 *
	 * @return array
	 */
	/*protected function getOptions()
	{
		return array(
			array('example', null, InputOption::VALUE_OPTIONAL, 'An example option.', null),
		);
	}*/

}
