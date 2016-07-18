from functools import wraps
from time import perf_counter
from SpaceDock.config import Config
from SpaceDock.routing import add_wrapper, route

# Load the cfg
cfg = Config('profiler.ini')

# Only register things if the profiler is enabled
if cfg.getb('profiler-histogram'):
    profiler_time = cfg.getf('profiler')
    histogram_data = {}

    def profile_method(f):
        @wraps(f)
        def profile_method_real(*args, **kwargs):
            profiler_name = f.__name__
            if hasattr(f, 'api_path'):
                profiler_name = f.api_path
            startTime = perf_counter()
            result = f(*args, **kwargs)
            endTime = perf_counter()
            timeDelta = 1000 * (endTime - startTime)
            if not profiler_time == 0 and timeDelta > profiler_time:
                print(profiler_name + " took " + str(timeDelta) + " ms")
            if not profiler_name in histogram_data:
                histogram_data[profiler_name] = {}
            timeDeltaInt = int(timeDelta)
            f_histogram_data = histogram_data[profiler_name]
            if not timeDeltaInt in f_histogram_data:
                f_histogram_data[timeDeltaInt] = 0
            f_histogram_data[timeDeltaInt] = f_histogram_data[timeDeltaInt] + 1
            return result
        return profile_method_real

    add_wrapper(profile_method)

    @route('/profiler/histogram')
    def histogram():
        return histogram_data