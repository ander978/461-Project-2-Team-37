using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using TrustRegistry.Models;
using TrustRegistry.Data;

namespace TrustRegistry.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class PackageController : ControllerBase
    {
        private readonly ApiContext _context;

        public PackageController(ApiContext context)
        {
            _context = context;
        }
        [HttpGet]
        public JsonResult GetRegistries()
        {
        var registries = _context.Registries.ToList();

        return new JsonResult(registries);
        }
        //Create and Edit 
        [HttpPost]
        public JsonResult CreateEdit(PackageRegistry registry)
        {
            if(registry.PackageRegistryId == 0)
            {
                _context.Registries.Add(registry);
            } else
            {
                var registryInDb = _context.Registries.Find(registry.PackageRegistryId);

                if (registryInDb == null)
                    return new JsonResult(NotFound());

                registryInDb = registry;
            }

            _context.SaveChanges();

            return new JsonResult(Ok(registry));

        }
    }
}
