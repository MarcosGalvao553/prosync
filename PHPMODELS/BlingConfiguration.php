<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class BlingConfiguration extends Model
{
    use HasFactory;

    protected $fillable = [
        'client_id',
        'secret_key',
        'url_callback',
        'postcode',
        'access_token',
        'refresh_token',
        'token_validate',
        'code',
        'user_id'
    ];    
}
