<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class ProductUser extends Model
{
    use HasFactory;

    protected $fillable = [
		'user_id',
		'product_id',
		'tiny_product_id',
		'bling_product_id',
		'price'
	];

    public function user()
	{
		return $this->belongsTo(User::class);
	}

    public function Product()
	{
		return $this->belongsTo(Product::class);
	}
}
