package utils

import (
	"math/rand"
	"time"

	"github.com/TinyKitten/TimelineServer/models"
	"gopkg.in/mgo.v2/bson"
)

func initUsers() (users []models.User) {
	users = append(users, *models.NewUser("52ba__", "password", "52ba@gmail.com", true))
	users[0].DisplayName = "ごつば"
	users = append(users, *models.NewUser("kotten", "password", "kotten@gmail.com", true))
	users[1].DisplayName = "こったん"
	users = append(users, *models.NewUser("empire_official", "password", "empire@gmail.com", true))
	users[2].DisplayName = "【公式】グンマー帝国"
	users = append(users, *models.NewUser("kitten_offical", "password", "kitten@gmail.com", true))
	users[3].DisplayName = "株式会社Kitten"
	users = append(users, *models.NewUser("kitton", "password", "kitton@gmail.com", true))
	users[4].DisplayName = "ホテルキットン"
	return
}

func initPosts() (posts []models.Post) {
	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "よっぽど古いお話なんで御座ございますよ。私の祖父じじいの子供の時分に居りました、「三さん」という猫なんで御座ございます。三毛みけだったんで御座ございますって。"))
	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "何でも、あの、その祖父じじいの話に、おばあさんがお嫁に来る時に――祖父じじいのお母さんなんで御座ございましょうねえ――泉州堺せんしゅうさかいから連れて来た猫なんで御座いますって。"))
	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "随分ずいぶん永く――家に十八年も居たんで御座ございますよ。大きくなっておりましたそうです。もう、耳なんか、厚ぼったく、五分ぶぐらいになっていたそうで御座ございますよ。"))

	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "ある夕がた、少年探偵団の名コンビ井上一郎いのうえいちろう君とノロちゃんとが、世田谷せたがや区のさびしいやしきまちを歩いていました。きょうは井上君のほうが"))
	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "ノロちゃんというのは、野呂一平のろいっぺい君のあだなです。ノロちゃんは団員のうちでいちばん、おくびょうものですが、ちゃめで、あいきょうもので、みんなにすかれています。"))
	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "「ねえ、ノロちゃん、ぼくたちは少年探偵団員だよ。このまま逃げだすわけにはいかない。あいつのあとをつけてみよう。お化けなんているはずがないよ。"))

	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "سارعي للمجد و العلياء"))
	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "مجدي لخالق السماء"))
	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "و ارفعي الخفاق أخضر"))

	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "כָּל עוֹד בַּלֵּבָב פְּנִימָה"))
	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "נֶפֶשׁ יְהוּדִי הוֹמִיָּה"))
	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "עַיִן לְצִיּוֹן צוֹפיָּה"))

	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "'eretz Tziyyôn wŶrūšāláyim"))
	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "lihəyōth ‘am chophšī bə'artzēnū,"))
	posts = append(posts, *models.NewPost(bson.NewObjectId(), bson.NewObjectId(), "‘ôdh lo' 'ābhdhāh tiqwāthēnū"))

	return
}

func genRand(max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}

func GenerateSamplePost() models.Post {
	posts := initPosts()
	return posts[genRand(len(posts))]
}

func GenerateSamplePostResponse() models.PostResponse {
	post := GenerateSamplePost()
	users := initUsers()
	return models.PostToPostResponse(post, users[genRand(len(users))])
}
