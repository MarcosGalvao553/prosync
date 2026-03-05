<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;

class PreSaleProducts extends Model
{
    protected $table = 'pre_sale_products';

    protected $fillable = [
        'product_id',
        'end_date',
        'active',
    ];

    protected $casts = [
        'end_date' => 'datetime',
        'active' => 'boolean',
    ];

    public function product()
    {
        return $this->belongsTo(Product::class);
    }
}
