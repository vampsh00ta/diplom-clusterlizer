// Глобальные переменные для состояния пагинации
let cursor = null;  // Глобальный курсор для пагинации
let hasNextPage = true;
let isLoading = false; // Флаг для предотвращения повторных запросов

// Функция для получения вакансий
async function fetchVacancies() {
    if (isLoading) {
        // Если уже идет загрузка, просто выходим из функции
        return [];
    }

    // Устанавливаем флаг загрузки
    isLoading = true;

    const search = document.getElementById('search').value;
    const city = document.getElementById('city').value;
    const experience = document.getElementById('experience').value;
    const company = document.getElementById('company').value;
    const speciality = document.getElementById('speciality').value;


    // const keyword_slugs = filters.join(" ")

    // const company = document.getElementById('company').value;

    // GraphQL-запрос
    const query = `
       query getVacancies(
       $cities: [String],
       $keyword_slugs: [String], 
       $experience_slug: String,  
       $speciality_slug: String,  
       $search: String,  
       $company_slug: String, 
       $first: Int, 
       $after: String)
        {
          vacancies(
              cities: $cities, 
              experience_slug: $experience_slug,
              speciality_slug: $speciality_slug,
              
              keyword_slugs: $keyword_slugs,  
              company_slug: $company_slug,  
              search: $search,
              first: $first, 
              after: $after) {
                edges {
                  node {
                    name
                    slug_id
                    cities
                    company {
                      name
                      slug
                    }
                    experience {
                      name
                      slug
                    }
                    category
                    link
                    description
                    
                  }
                  cursor
                }
                pageInfo {
                  hasNextPage
                  endCursor
                }
              }
        }
    `;

    // Данные для GraphQL запроса
    description = ""
    if (search.length !=0){
        description+=search
    }
    // if (slugsStr.length!=0){
    //     description+=slugsStr
    // }
    const variables = {
        search: search,
        cities: city ? [city] : [],
        experience_slug: experience,
        company_slug: company,
        speciality_slug: speciality,

        keyword_slugs: filters,
        first: 20,
        after: cursor, // Используем текущий курсор для пагинации
    };

    try {
        // Запрос к API
        const response = await fetch('/graphql', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                query: query,
                variables: variables
            })
        });

        const result = await response.json();

        const vacanciesData = result.data.vacancies;

        // Обновляем глобальный курсор и состояние пагинации
        cursor = vacanciesData.pageInfo.endCursor; // Обновляем курсор на значение из последнего запроса
        hasNextPage = vacanciesData.pageInfo.hasNextPage; // Обновляем информацию о наличии следующей страницы



        return vacanciesData.edges.map(edge => edge.node);
    } catch (error) {

        return [];
    } finally {
        // Сбрасываем флаг загрузки после завершения запроса
        isLoading = false;
    }
}

// Обработчик нажатия на кнопку "Найти"
document.getElementById('search-button').addEventListener('click', async () => {
    // Обнуляем курсор и флаг наличия следующей страницы, так как фильтры изменились
    cursor = null;
    hasNextPage = true;

    // Очищаем контейнер с результатами, чтобы убрать предыдущие вакансии
    const resultsContainer = document.getElementById('results');
    resultsContainer.innerHTML = '';

    // Показываем loader во время выполнения запроса
    document.getElementById('loader').classList.remove('hidden');

    // Выполняем новый запрос с учетом измененных фильтров
    const vacancies = await fetchVacancies();

    // Обновляем результаты на странице
    updateResults(vacancies);

    // Скрываем loader после завершения
    document.getElementById('loader').classList.add('hidden');
});

// Обработчик прокрутки страницы для подгрузки новых вакансий
window.addEventListener('scroll', async () => {
    // Проверяем, достиг ли пользователь конца страницы и не идет ли загрузка
    if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight - 100 && hasNextPage && !isLoading) {
        console.log("Начинается загрузка новых вакансий..."); // Отладка
        document.getElementById('loader').classList.remove('hidden'); // Показываем loader перед началом загрузки
        const vacancies = await fetchVacancies(); // Используем обновленный курсор для следующего запроса
        updateResults(vacancies);
        document.getElementById('loader').classList.add('hidden'); // Скрываем loader после завершения загрузки
    }
});

function updateResults(vacancies) {
    const resultsContainer = document.getElementById('results');
    const vacancyTemplate = document.getElementById('vacancy-template').content;

    // Если есть вакансии
    if (vacancies && vacancies.length > 0) {
        vacancies.forEach(vacancy => {
            console.log("Добавление вакансии:", vacancy); // <-- Отладка добавляемой вакансии
            const clone = vacancyTemplate.cloneNode(true);
            clone.querySelector('.vacancy-title').textContent = vacancy.name;
            clone.querySelector('.vacancy-link').href = vacancy.link;
            clone.querySelector('.vacancy-company').textContent = `Компания: ${vacancy.company.name}`;
            clone.querySelector('.vacancy-experience').textContent = `Опыт работы: ${vacancy.experience.name}`;
            clone.querySelector('.vacancy-city').textContent = `Город: ${vacancy.cities[0]}`;
            resultsContainer.appendChild(clone);
        });
    } else {
        // Если нет вакансий, показываем сообщение
        console.log("Вакансии не найдены для нового поиска");
        resultsContainer.innerHTML = '<p class="text-center text-gray-500">Вакансии не найдены</p>';
    }
}