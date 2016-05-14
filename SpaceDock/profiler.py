from flask import jsonify
from functools import wraps
from time import perf_counter

class Profiler:
    def __init__(self, cfg):
        self.cfg = cfg
        self.profiler_time = self.cfg.getf('profiler')
        self.profiler_histogram = self.cfg.getb('profiler-histogram')
        self.histogram_data = {}
        
    def profile_method(self, f):
        @wraps(f)
        def profile_method_real(*args, **kwargs):
            profiler_name = f.__name__
            if hasattr(f, 'api_path'):
                profiler_name = f.api_path
            startTime = perf_counter()
            result = f(*args, **kwargs)
            endTime = perf_counter()
            timeDelta = 1000 * (endTime - startTime)
            if not self.profiler_time == 0 and timeDelta > self.profiler_time:
                print(profiler_name + " took " + str(timeDelta) + " ms")
            if self.profiler_histogram:      
                if not profiler_name in self.histogram_data:
                    self.histogram_data[profiler_name] = {}
                timeDeltaInt = int(timeDelta)
                f_histogram_data = self.histogram_data[profiler_name]
                if not timeDeltaInt in f_histogram_data:
                    f_histogram_data[timeDeltaInt] = 0
                f_histogram_data[timeDeltaInt] = f_histogram_data[timeDeltaInt] + 1
            return result
        return profile_method_real
    
    def histogram(self):
        return jsonify(self.histogram_data)
        
    histogram.api_path = "/profiler/histogram"