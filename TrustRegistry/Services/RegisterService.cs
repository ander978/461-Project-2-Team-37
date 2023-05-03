using System;
using System.Collections.Generic;
using System.Security.Cryptography;
using System.Text;
using System.Threading.Tasks;
using MongoDB.Bson;
using MongoDB.Driver;
using Newtonsoft.Json;

using TrustRegistry.Models;

namespace TrustRegistry.Services
{
    public class RegisterService
    {
        private readonly string _usersCollectionName = "Users";
        private readonly string _tokenUsageCollectionName = "TokenUsage";
        private readonly string _jwtSecret = "GIyoqsMwGPv2YEStDNat1qaXXbOH8lmwkvbUODyzoF8="; // Replace with your own secret
        private readonly long _expirationTime = 36000000; // 1 day in milliseconds
        private readonly long _maxTokenUsage = 1000;

        private readonly IMongoCollection<User> _usersCollection;
        private readonly IMongoCollection<TokenUsage> _tokenUsageCollection;

        public RegisterService(IMongoDatabase database)
        {
            _usersCollection = database.GetCollection<User>(_usersCollectionName);
            _tokenUsageCollection = database.GetCollection<TokenUsage>(_tokenUsageCollectionName);
        }

        public async Task<User> SaveUser(User user, UserAuthInf authInfo)
        {
            var existingUser = await _usersCollection.Find(u => u.Name == user.Name).FirstOrDefaultAsync();

            if (existingUser != null)
            {
                return null; // User already exists
            }

            user.UserAuthInf = new Dictionary<string, string>
            {
                { "password", authInfo.Password }
            };

            await _usersCollection.InsertOneAsync(user);

            return user;
        }

        public async Task RemoveUser(string username)
        {
            var result = await _usersCollection.DeleteOneAsync(u => u.Name == username);

            if (result.DeletedCount == 0)
            {
                throw new ArgumentException("User not found.");
            }
        }

        public async Task<string> AuthenticateUser(User user, UserAuthInf authInfo)
        {
            var dbUser = await _usersCollection.Find(u => u.Name == user.Name).FirstOrDefaultAsync();

            if (dbUser == null)
            {
                throw new ArgumentException("User not found.");
            }

            var storedPassword = dbUser.UserAuthenticationInfo["password"];

            if (authInfo.Password != storedPassword)
            {
                throw new ArgumentException("Invalid password.");
            }

            var token = GenerateJwtToken(user.Name);
            await StoreTokenUsage(new TokenUsage { Token = token, Username = user.Name });

            return token;
        }

        private string GenerateJwtToken(string username)
        {
            var algorithm = new HMACSHA256(Encoding.UTF8.GetBytes(_jwtSecret));
            var now = DateTime.UtcNow;
            var expirationDate = now.AddMilliseconds(_expirationTime);

            var payload = new Dictionary<string, object>
            {
                { "iss", "auth0" },
                { "sub", username },
                { "iat", now.ToUnixTimeSeconds() },
                { "exp", expirationDate.ToUnixTimeSeconds() },
            };

            return JWT.Encode(payload, algorithm, JwsAlgorithm.HS256);
        }
    }
}
    