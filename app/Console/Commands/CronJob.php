<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;

class CronJob extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'concierge:cronjob';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Sets up the cronjobs for concierge to remove expired/obsolete leases.';

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
    public function handle()
    {
        //Get the current crontabs
        $output = shell_exec('crontab -l');
        $cron_time = env('CRON_TIME', '* * * * *');
        //The cron needed for this repo
        $cron = $cron_time . ' ' . config('concierge.php_path') . ' '
            . config('concierge.artisan_path') . ' concierge:clean-lease';

        //Check if Cron Already exists
        if (strpos($output, $cron) === false) {
            //setup Old Cron with the new cron included
            file_put_contents('/tmp/crontab.txt', $output.$cron.PHP_EOL);
            exec('crontab /tmp/crontab.txt');
            $this->info('Cron Added for running at ' . $cron_time);
        } else {
            $this->info('Cron Already exists');
        }
    }
}
