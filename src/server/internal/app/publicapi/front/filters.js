const filtersContainer = document.getElementById('filters-container');
const filtersInput = document.getElementById('filters-input');
const filters = [];
availableKeywords = [];

// Обработчик нажатия на кнопку "Сохранить фильтры"
document.getElementById('save-filters-button').addEventListener('click', async () => {
    const user = checkAuthorization();
    if (!user) return; // Если не авторизован, прекратить выполнение

    // Получаем значения фильтров
    const city = document.getElementById('city').value;
    const experience = document.getElementById('experience').value;
    const company = document.getElementById('company').value;
    // 1. Получаем выбранную специализацию
    const speciality = document.getElementById('speciality').value;

    // Собираем слуги (теги)
    const keywords = filters;

    // Формируем данные для сохранения
    const saveFilterData = {
        tg_id: parseInt(user.id, 10),
        city: city,
        experience: experience,
        company: company,
        speciality: speciality, // 2. Добавляем поле speciality
        keywords: keywords
    };

    try {
        const response = await fetch('/api/v1/user/saveFilter', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': Object.entries(user).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`).join('&')
            },
            body: JSON.stringify(saveFilterData)
        });

        if (!response.ok) {
            throw new Error('Ошибка при сохранении фильтров');
        }

        const responseData = await response.json();
        console.log('Фильтры сохранены:', responseData);
        alert('Фильтры успешно сохранены!');

    } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка при сохранении фильтров');
    }
});

// Функция загрузки данных всех фильтров (города, компании, опыт работы, специальности, ключевые слова)
async function loadAllFilters() {
    const query = `
        query {
            allFilters {
                experiences {
                    slug
                    name
                }
                cities {
                    name
                    vacancy_count
                }
                companies {
                    name
                    slug
                }
                keywords {
                    name
                    slug
                }
                specialities {
                    name
                    slug
                }
            }
        }
    `;

    const response = await fetch('/graphql', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ query })
    });

    const result = await response.json();
    return result.data;
}

// Функция загрузки пользовательских фильтров
async function loadUserFilters() {
    const user = checkAuthorization();
    if (!user) return; // Если не авторизован, прекратить выполнение

    try {
        const response = await fetch(`/api/v1/user/filter?tg_id=${user.id}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': Object.entries(user).map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`).join('&')
            }
        });

        if (!response.ok) {
            throw new Error('Ошибка при получении фильтров');
        }

        const data = await response.json();

        // Устанавливаем значения города
        document.getElementById('city').value = data.city || "Любой";

        // Опыт работы
        const experienceSelect = document.getElementById('experience');
        const experienceOptions = Array.from(experienceSelect.options);
        const userExperience = data.experience?.slug || "";
        const matchedExperience = experienceOptions.find(option => option.value === userExperience);
        if (matchedExperience) {
            matchedExperience.selected = true;
        } else {
            experienceSelect.value = ""; // "Любой", если пользовательский опыт не найден
        }

        // Компания
        const companySelect = document.getElementById('company');
        const companyOptions = Array.from(companySelect.options);
        const userCompany = data.company?.slug || "";
        const matchedCompany = companyOptions.find(option => option.value === userCompany);
        if (matchedCompany) {
            matchedCompany.selected = true;
        } else {
            companySelect.value = ""; // "Любая", если компания не найдена
        }

        // 3. Специализация
        const specialitySelect = document.getElementById('speciality');
        const specialityOptions = Array.from(specialitySelect.options);
        const userSpeciality = data.speciality?.slug || "";
        const matchedSpeciality = specialityOptions.find(option => option.value === userSpeciality);
        if (matchedSpeciality) {
            matchedSpeciality.selected = true;
        } else {
            specialitySelect.value = ""; // "Любая", если специализация не найдена
        }

        // Дополнительные фильтры (keywords)
        if (data.keywords && data.keywords.length > 0) {
            data.keywords.forEach(keyword => {
                const availableKeyword = availableKeywords.find(kw => kw.slug === keyword.slug);
                if (availableKeyword) {
                    addFilter(availableKeyword.slug, availableKeyword.name);
                }
            });
        }

    } catch (error) {
        console.error('Ошибка:', error);
    }
}

// Обработчик загрузки данных фильтров после загрузки страницы
document.addEventListener('DOMContentLoaded', async () => {
    const data = await loadAllFilters();
    const allFiltersData = data.allFilters;

    // Загрузка ключевых слов (для дополнительных фильтров)
    availableKeywords = allFiltersData.keywords.map(keyword => ({
        name: keyword.name,
        slug: keyword.slug
    }));

    // Города
    const cities = allFiltersData.cities;
    const citySelect = document.getElementById('city');
    citySelect.innerHTML = '<option value="">Любой</option>'; // Опция по умолчанию
    cities.forEach(city => {
        const option = document.createElement('option');
        option.value = city.name;
        option.textContent = `${city.name} (${city.vacancy_count} вакансий)`;
        citySelect.appendChild(option);
    });

    // Компании
    const companies = allFiltersData.companies;
    const companySelect = document.getElementById('company');
    companySelect.innerHTML = '<option value="">Любая</option>'; // Опция по умолчанию
    companies.forEach(company => {
        const option = document.createElement('option');
        option.value = company.slug;
        option.textContent = company.name;
        companySelect.appendChild(option);
    });

    // Опыт работы
    const experiences = allFiltersData.experiences;
    const experienceSelect = document.getElementById('experience');
    experienceSelect.innerHTML = '<option value="">Любой</option>'; // Опция по умолчанию
    experiences.forEach(experience => {
        const option = document.createElement('option');
        option.value = experience.slug;
        option.textContent = experience.name;
        experienceSelect.appendChild(option);
    });

    // 4. Загрузка специальностей
    const specialities = allFiltersData.specialities;
    const specialitySelect = document.getElementById('speciality');
    specialitySelect.innerHTML = '<option value="">Любая</option>';
    specialities.forEach(speciality => {
        const option = document.createElement('option');
        option.value = speciality.slug;
        option.textContent = speciality.name;
        specialitySelect.appendChild(option);
    });

    // Загружаем пользовательские фильтры (если есть)
    loadUserFilters();
});

// Добавление фильтра при вводе пробела
filtersInput.addEventListener('keyup', (event) => {
    if (event.key === ' ' && filtersInput.value.trim() !== '') {
        const filterText = filtersInput.value.trim();
        addFilter(filterText, filterText);
        filtersInput.value = ''; // Очищаем поле ввода
    }
});

// Функция добавления фильтра
function addFilter(slug, name) {
    if (!filters.includes(slug)) {
        filters.push(slug);
        const filterElement = document.createElement('div');
        filterElement.classList.add(
            'filter', 'flex', 'items-center', 'bg-blue-500', 'text-white',
            'px-2', 'py-1', 'rounded-full', 'mr-2', 'mb-2'
        );
        filterElement.innerHTML = `
            <span class="mr-1">${name}</span>
            <button class="remove-filter focus:outline-none" title="Удалить фильтр">&times;</button>
        `;

        // Добавляем фильтр внутрь контейнера, находящегося рядом с полем ввода
        filtersContainer.insertBefore(filterElement, filtersInput);

        // Обработчик удаления фильтра
        filterElement.querySelector('.remove-filter').addEventListener('click', () => {
            removeFilter(slug, filterElement);
        });
    }
}

// Функция удаления фильтра
function removeFilter(slug, element) {
    filters.splice(filters.indexOf(slug), 1);
    filtersContainer.removeChild(element);
}

// Обработчик нажатия на кнопку "Фильтры"
document.getElementById('filters-toggle-button').addEventListener('click', () => {
    const filtersMenu = document.getElementById('filters-menu');
    if (filtersMenu.classList.contains('hidden')) {
        filtersMenu.classList.remove('hidden');
    } else {
        filtersMenu.classList.add('hidden');
    }
});

// Обработчик ввода текста для дополнительных фильтров (подсказки)
filtersInput.addEventListener('input', () => {
    const inputValue = filtersInput.value.toLowerCase();
    const suggestions = availableKeywords.filter(keyword =>
        keyword.name.toLowerCase().includes(inputValue)
    );

    // Удаляем старые подсказки
    const oldSuggestionBox = document.getElementById('suggestion-box');
    if (oldSuggestionBox) {
        oldSuggestionBox.remove();
    }

    // Создаем контейнер для подсказок
    const suggestionBox = document.createElement('div');
    suggestionBox.id = 'suggestion-box';
    suggestionBox.classList.add('suggestion-box');
    suggestionBox.style.left = `${filtersInput.offsetLeft}px`;
    suggestionBox.style.top = `${filtersInput.offsetTop + filtersInput.offsetHeight}px`;
    suggestionBox.style.width = `${filtersInput.offsetWidth}px`;

    // Наполняем подсказками
    suggestions.forEach(suggestion => {
        const suggestionItem = document.createElement('div');
        suggestionItem.classList.add('suggestion-item');
        suggestionItem.textContent = suggestion.name;
        suggestionItem.dataset.slug = suggestion.slug;

        // При клике на подсказку добавляем фильтр
        suggestionItem.addEventListener('click', () => {
            addFilter(suggestionItem.dataset.slug, suggestionItem.textContent);
            filtersInput.value = '';
            suggestionBox.remove();
        });

        suggestionBox.appendChild(suggestionItem);
    });

    filtersContainer.appendChild(suggestionBox);
});
