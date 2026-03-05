<?php

/**
 * Created by Reliese Model.
 */

namespace App\Models;

use Carbon\Carbon;
use Illuminate\Database\Eloquent\Collection;
use Illuminate\Database\Eloquent\Model;

/**
 * Class Product
 * 
 * @property int $id
 * @property string $name
 * @property string|null $description
 * @property float|null $price
 * @property float $cost_price
 * @property int $category_id
 * @property bool $isEnabled
 * @property int $sale_count
 * @property int $review_count
 * @property string|null $sku
 * @property int|null $stock
 * @property string|null $observation
 * @property float|null $weight
 * @property float|null $height
 * @property float|null $width
 * @property float|null $length
 * @property string|null $product_tiny
 * 
 * @property Category $category
 * @property Collection|ProductImage[] $product_images
 * @property Collection|ProductOrderItem[] $product_order_items
 * @property Collection|PropertyValue[] $property_values
 * @property Collection|Review[] $reviews
 * @property Collection|StoreFeatured[] $store_featureds
 * @property Collection|StoreShippingOrderItem[] $store_shipping_order_items
 *
 * @package App\Models
 */
class Product extends Model
{
	protected $table = 'products';
	// public $timestamps = false;

	protected $casts = [
		'price' => 'float',
		'cost_price' => 'float',
		'category_id' => 'int',
		'isEnabled' => 'bool',
		'isPreSale' => 'bool',
		'sale_count' => 'int',
		'review_count' => 'int',
		'stock' => 'int',
		'weight' => 'float',
		'height' => 'float',
		'width' => 'float',
		'length' => 'float'
	];

	protected $fillable = [
		'name',
		'description',
		'price',
		'cost_price',
		'category_id',
		'isEnabled',
		'isPreSale',
		'sale_count',
		'review_count',
		'sku',
		'stock',
		'observation',
		'weight',
		'height',
		'width',
		'length',
		'product_tiny',
		'ncm',
		'ean',
		'marca',
		'cest',
		'stop_stock',
		'promotion_id',
		'original_price'
	];

	public function category()
	{
		return $this->belongsTo(Category::class);
	}

	public function product_images()
	{
		return $this->hasMany(ProductImage::class);
	}

	public function product_users()
	{
		return $this->hasMany(ProductUser::class);
	}

	public function product_order_items()
	{
		return $this->hasMany(ProductOrderItem::class);
	}

	public function property_values()
	{
		return $this->hasMany(PropertyValue::class);
	}

	public function reviews()
	{
		return $this->hasMany(Review::class);
	}

	public function store_featureds()
	{
		return $this->hasMany(StoreFeatured::class);
	}

	public function store_shipping_order_items()
	{
		return $this->hasMany(StoreShippingOrderItem::class);
	}

	public function product_promotions()
	{
		return $this->hasMany(ProductPromotions::class);
	}

	public function product_pre_sale()
	{
		return $this->hasOne(PreSaleProducts::class);
	}
}
