<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use App\Models\Lease;

class LeaseManager extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'concierge:clean-lease';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Manages termination & removal of expired/obsolete leases.';

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
        $result = Lease::cleanLeases();
        if ($result) {
            $this->info($result);
        }
    }
}
