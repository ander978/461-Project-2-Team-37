using Microsoft.AspNetCore.Mvc;
using System;
using System.Threading.Tasks;
using System.Text;
using System.Security.Claims;
using Microsoft.IdentityModel.Tokens;
using System.IdentityModel.Tokens.Jwt;
using Microsoft.Extensions.Configuration;
//using com.spring_rest_api.api_paths.entity;
using TrustRegistry.Models;
using TrustRegistry.Services;
//using com.spring_rest_api.api_paths.service;

namespace TrustRegistry.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class RegisterController : ControllerBase
    {
        private readonly IConfiguration _config;
        private readonly RegisterService _registerService;

        public RegisterController(IConfiguration config, RegisterService registerService)
        {
            _config = config;
            _registerService = registerService;
        }

        [HttpPost("register")]
        public async Task<IActionResult> Register(AuthRequest request)
        {
            try
            {
                User registeredUser = await _registerService.SaveUser(request.User, request.UserAuthInf);
                if (registeredUser != null)
                {
                    return Ok(registeredUser);
                }
                else
                {
                    return BadRequest("User registration failed.");
                }
            }
            catch (Exception e)
            {
                return StatusCode(500, "Error registering the user.");
            }
        }

        [HttpPost("authenticate")]
        public async Task<IActionResult> Authenticate(AuthRequest request)
        {
            var user = await _registerService.AuthenticateUser(request.User, request.UserAuthInf);
            if (user == null)
            {
                return Unauthorized("Invalid username or password.");
            }

            var tokenString = GenerateJSONWebToken(user);
            return Ok(new { token = tokenString });
        }

        private string GenerateJSONWebToken(User user)
        {
            var securityKey = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(_config["Jwt:Key"]));
            var credentials = new SigningCredentials(securityKey, SecurityAlgorithms.HmacSha256);

            var claims = new[] {
                new Claim(JwtRegisteredClaimNames.Sub, user.Id.ToString()),
                new Claim(JwtRegisteredClaimNames.Email, user.Email),
                new Claim("role", user.Role)
            };

            var token = new JwtSecurityToken(
                _config["Jwt:Issuer"],
                _config["Jwt:Issuer"],
                claims,
                expires: DateTime.Now.AddMinutes(120),
                signingCredentials: credentials);

            return new JwtSecurityTokenHandler().WriteToken(token);
        }

        [HttpDelete("{userId}")]
        public async Task<IActionResult> RemoveUser(string userId)
        {
            try
            {
                await _registerService.RemoveUser(userId);
                return Ok("User successfully removed.");
            }
            catch (Exception e)
            {
                return StatusCode(500, "Error removing the user.");
            }
        }
    }
}