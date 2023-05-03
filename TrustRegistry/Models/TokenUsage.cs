using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace TrustRegistry.Models
{
    public class TokenUsage
    {
        [BsonId]
        [BsonRepresentation(BsonType.ObjectId)]
        public string Id { get; set; }

        [BsonElement("token")]
        public string Token { get; set; }

        [BsonElement("username")]
        public string Username { get; set; }

        [BsonElement("usage_count")]
        public long UsageCount { get; set; }
    }
}