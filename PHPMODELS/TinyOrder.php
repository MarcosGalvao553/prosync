<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class TinyOrder extends Model
{
    use HasFactory;

    protected $fillable = [
		'shipping_order_id',
        'order_tiny_id',
	];

}
