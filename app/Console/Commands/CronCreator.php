<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Input\InputArgument;

class CronCreator extends Command {

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
    protected $description = 'Sets up the cronjobs for concierge';

    /**
     * Execute the console command.
     *
     * @return mixed
     */
    public function handle()
    {
        //Get the current crontabs
        $output = shell_exec('crontab -l');

        //The cron needed for this repo
        $cron="* * * * *  ".config('concierge.php_path')." ".config('concierge.artisan_path')." custom:leasemanager";

        //Check if Cron Already exists
        if(strpos($output,$cron)===FALSE)
        {
            //setup Old Cron with the new cron included
            file_put_contents('/tmp/crontab.txt', $output.$cron.PHP_EOL);
            exec('crontab /tmp/crontab.txt');
            $this->info("Cron Added for running every minute");
        }
        else
        {
            $this->info("Cron Already exists");
        }

    }
}
